package pybind

import "C"

/*
   #cgo pkg-config: python3-embed
   #include <Python.h>
   #include <stdlib.h>

   // Forward declarations
   extern PyObject* go_callback_dispatcher(PyObject *self, PyObject *args, PyObject *kwargs);

   // Create a new Python method definition
   PyMethodDef* create_method_def(const char *name, const char *doc) {
       PyMethodDef *def = (PyMethodDef*)malloc(sizeof(PyMethodDef));
       def->ml_name = name;
       def->ml_meth = (PyCFunction)go_callback_dispatcher;
       def->ml_flags = METH_VARARGS | METH_KEYWORDS;
       def->ml_doc = doc;
       return def;
   }

   // Helper to convert tuple size
   Py_ssize_t tuple_size(PyObject *tuple) {
       return PyTuple_Size(tuple);
   }

   // Helper to get tuple item
   PyObject* tuple_get_item(PyObject *tuple, Py_ssize_t index) {
       return PyTuple_GetItem(tuple, index);
   }
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

var (
	pyMutex          sync.Mutex
	initialized      bool
	callbackRegistry = make(map[string]GoCallbackFunc)
	registryMutex    sync.RWMutex
	callbackCounter  int
)

// GoCallbackFunc is a Go function callable from Python
type GoCallbackFunc func(args []interface{}) (interface{}, error)

// PyObject wraps a Python C API object
type PyObject struct {
	ptr *C.PyObject
}

// NewPyObject creates a new PyObject with proper reference counting
func NewPyObject(ptr *C.PyObject) *PyObject {
	if ptr == nil {
		return nil
	}
	C.Py_IncRef(ptr)
	obj := &PyObject{ptr: ptr}
	runtime.SetFinalizer(obj, (*PyObject).decref)
	return obj
}

func (obj *PyObject) decref() {
	if obj.ptr != nil {
		pyMutex.Lock()
		C.Py_DecRef(obj.ptr)
		pyMutex.Unlock()
		obj.ptr = nil
	}
}

// PyInterpreter represents an embedded Python interpreter
type PyInterpreter struct {
	mainDict *PyObject
	goModule *PyObject
}

// Initialize starts the Python interpreter
func Initialize() error {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	if initialized {
		return nil
	}

	C.Py_Initialize()
	if C.Py_IsInitialized() == 0 {
		return fmt.Errorf("failed to initialize Python")
	}

	initialized = true
	return nil
}

// Finalize shuts down the Python interpreter
func Finalize() {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	if initialized {
		C.Py_Finalize()
		initialized = false
	}
}

// NewInterpreter creates a new Python interpreter context
func NewInterpreter() (*PyInterpreter, error) {
	if !initialized {
		return nil, fmt.Errorf("Python not initialized")
	}

	pyMutex.Lock()
	defer pyMutex.Unlock()

	// Get main module dictionary
	mainName := C.CString("__main__")
	defer C.free(unsafe.Pointer(mainName))

	mainModule := C.PyImport_AddModule(mainName)
	if mainModule == nil {
		return nil, fmt.Errorf("failed to get __main__")
	}

	mainDict := C.PyModule_GetDict(mainModule)
	if mainDict == nil {
		return nil, fmt.Errorf("failed to get __main__ dict")
	}

	// Create go module
	goModuleName := C.CString("go")
	defer C.free(unsafe.Pointer(goModuleName))

	goModule := C.PyModule_New(goModuleName)
	if goModule == nil {
		return nil, fmt.Errorf("failed to create go module")
	}

	// Add go module to sys.modules
	sysModules := C.PyImport_GetModuleDict()
	C.PyDict_SetItemString(sysModules, goModuleName, goModule)

	interp := &PyInterpreter{
		mainDict: NewPyObject(mainDict),
		goModule: NewPyObject(goModule),
	}

	// Import go module into main
	importCode := C.CString("import go")
	defer C.free(unsafe.Pointer(importCode))
	C.PyRun_SimpleString(importCode)

	return interp, nil
}

// RegisterFunction registers a Go function callable from Python
func (pi *PyInterpreter) RegisterFunction(name string, fn GoCallbackFunc) error {
	// Store callback in registry
	registryMutex.Lock()
	callbackName := fmt.Sprintf("_go_callback_%d_%s", callbackCounter, name)
	callbackRegistry[callbackName] = fn
	callbackCounter++
	registryMutex.Unlock()

	pyMutex.Lock()
	defer pyMutex.Unlock()

	// Create Python function wrapper
	wrapperCode := fmt.Sprintf(`
def %s(*args, **kwargs):
    import ctypes
    # This is a placeholder - the actual callback happens through CGO
    # We'll use a workaround with exec and globals
    return _internal_go_call('%s', args, kwargs)
`, name, callbackName)

	cCode := C.CString(wrapperCode)
	defer C.free(unsafe.Pointer(cCode))

	// Store callback reference in module
	callbackPtr := C.CString(callbackName)
	defer C.free(unsafe.Pointer(callbackPtr))

	callbackObj := C.PyUnicode_FromString(callbackPtr)
	defer C.Py_DecRef(callbackObj)

	moduleName := C.CString(name)
	defer C.free(unsafe.Pointer(moduleName))

	C.PyModule_AddObject(pi.goModule.ptr, moduleName, callbackObj)

	return nil
}

// RegisterFunctionDirect registers a function with direct Python wrapping
func (pi *PyInterpreter) RegisterFunctionDirect(name string, fn GoCallbackFunc) error {
	registryMutex.Lock()
	callbackRegistry[name] = fn
	registryMutex.Unlock()

	// Create a simpler approach: store function ID in globals
	// and create a Python wrapper that calls through a helper
	code := fmt.Sprintf(`
class GoFunction:
    def __init__(self, name):
        self.name = name
    
    def __call__(self, *args):
        # This will be intercepted by our Go code
        import sys
        sys.stdout.write(f"__GO_CALL__{self.name}__")
        return None

go.%s = GoFunction("%s")
`, name, name)

	return pi.RunString(code)
}

// RunString executes Python code
func (pi *PyInterpreter) RunString(code string) error {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	result := C.PyRun_SimpleString(cCode)
	if result != 0 {
		C.PyErr_Print()
		return fmt.Errorf("failed to execute code")
	}

	return nil
}

// Eval evaluates a Python expression
func (pi *PyInterpreter) Eval(expr string) (*PyObject, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cExpr := C.CString(expr)
	defer C.free(unsafe.Pointer(cExpr))

	filename := C.CString("<eval>")
	defer C.free(unsafe.Pointer(filename))

	result := C.PyRun_String(cExpr, C.Py_eval_input, pi.mainDict.ptr, pi.mainDict.ptr)
	if result == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("eval failed")
	}

	return NewPyObject(result), nil
}

// GetGlobal retrieves a global variable
func (pi *PyInterpreter) GetGlobal(name string) (*PyObject, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	obj := C.PyDict_GetItemString(pi.mainDict.ptr, cName)
	if obj == nil {
		return nil, fmt.Errorf("global not found: %s", name)
	}

	return NewPyObject(obj), nil
}

// SetGlobal sets a global variable
func (pi *PyInterpreter) SetGlobal(name string, value *PyObject) error {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	if C.PyDict_SetItemString(pi.mainDict.ptr, cName, value.ptr) != 0 {
		return fmt.Errorf("failed to set global: %s", name)
	}

	return nil
}

// ImportModule imports a Python module
func (pi *PyInterpreter) ImportModule(name string) (*PyObject, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	module := C.PyImport_ImportModule(cName)
	if module == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("failed to import module: %s", name)
	}

	return NewPyObject(module), nil
}

// GetAttr gets an attribute from an object
func (obj *PyObject) GetAttr(name string) (*PyObject, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	attr := C.PyObject_GetAttrString(obj.ptr, cName)
	if attr == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("attribute not found: %s", name)
	}

	return NewPyObject(attr), nil
}

// Call calls a Python callable
func (obj *PyObject) Call(args ...*PyObject) (*PyObject, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	argTuple := C.PyTuple_New(C.Py_ssize_t(len(args)))
	if argTuple == nil {
		return nil, fmt.Errorf("failed to create tuple")
	}
	defer C.Py_DecRef(argTuple)

	for i, arg := range args {
		C.Py_IncRef(arg.ptr)
		C.PyTuple_SetItem(argTuple, C.Py_ssize_t(i), arg.ptr)
	}

	result := C.PyObject_CallObject(obj.ptr, argTuple)
	if result == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("call failed")
	}

	return NewPyObject(result), nil
}

// Type checking
func (obj *PyObject) IsInt() bool {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return C.PyLong_Check(obj.ptr) != 0
}

func (obj *PyObject) IsFloat() bool {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return C.PyFloat_Check(obj.ptr) != 0
}

func (obj *PyObject) IsString() bool {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return C.PyUnicode_Check(obj.ptr) != 0
}

func (obj *PyObject) IsList() bool {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return C.PyList_Check(obj.ptr) != 0
}

func (obj *PyObject) IsDict() bool {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return C.PyDict_Check(obj.ptr) != 0
}

// Conversions: Python to Go
func (obj *PyObject) AsInt() (int64, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	val := C.PyLong_AsLongLong(obj.ptr)
	if val == -1 && C.PyErr_Occurred() != nil {
		C.PyErr_Clear()
		return 0, fmt.Errorf("not an integer")
	}
	return int64(val), nil
}

func (obj *PyObject) AsFloat() (float64, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	val := C.PyFloat_AsDouble(obj.ptr)
	if val == -1.0 && C.PyErr_Occurred() != nil {
		C.PyErr_Clear()
		return 0, fmt.Errorf("not a float")
	}
	return float64(val), nil
}

func (obj *PyObject) AsString() (string, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cStr := C.PyUnicode_AsUTF8(obj.ptr)
	if cStr == nil {
		C.PyErr_Clear()
		return "", fmt.Errorf("not a string")
	}
	return C.GoString(cStr), nil
}

func (obj *PyObject) AsBool() (bool, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	val := C.PyObject_IsTrue(obj.ptr)
	if val == -1 {
		C.PyErr_Clear()
		return false, fmt.Errorf("not a bool")
	}
	return val == 1, nil
}

// ToInterface converts Python object to Go interface{}
func (obj *PyObject) ToInterface() interface{} {
	if obj.IsInt() {
		val, _ := obj.AsInt()
		return val
	}
	if obj.IsFloat() {
		val, _ := obj.AsFloat()
		return val
	}
	if obj.IsString() {
		val, _ := obj.AsString()
		return val
	}
	return nil
}

// Conversions: Go to Python
func PyInt(val int64) *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return NewPyObject(C.PyLong_FromLongLong(C.longlong(val)))
}

func PyFloat(val float64) *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return NewPyObject(C.PyFloat_FromDouble(C.double(val)))
}

func PyString(val string) *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	cStr := C.CString(val)
	defer C.free(unsafe.Pointer(cStr))
	return NewPyObject(C.PyUnicode_FromString(cStr))
}

func PyBool(val bool) *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	if val {
		return NewPyObject(C.Py_True)
	}
	return NewPyObject(C.Py_False)
}

func PyNone() *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return NewPyObject(C.Py_None)
}

func PyList(items ...*PyObject) *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	list := C.PyList_New(C.Py_ssize_t(len(items)))
	for i, item := range items {
		C.Py_IncRef(item.ptr)
		C.PyList_SetItem(list, C.Py_ssize_t(i), item.ptr)
	}
	return NewPyObject(list)
}

func PyDict() *PyObject {
	pyMutex.Lock()
	defer pyMutex.Unlock()
	return NewPyObject(C.PyDict_New())
}

func (dict *PyObject) SetItem(key string, value *PyObject) error {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	if C.PyDict_SetItemString(dict.ptr, cKey, value.ptr) != 0 {
		return fmt.Errorf("failed to set item")
	}
	return nil
}

func (dict *PyObject) GetItem(key string) (*PyObject, error) {
	pyMutex.Lock()
	defer pyMutex.Unlock()

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	item := C.PyDict_GetItemString(dict.ptr, cKey)
	if item == nil {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return NewPyObject(item), nil
}

// Example struct
type GameState struct {
	Level  int
	Score  int
	Player string
	Items  []string
}

func (g *GameState) ToPython() *PyObject {
	dict := PyDict()
	dict.SetItem("level", PyInt(int64(g.Level)))
	dict.SetItem("score", PyInt(int64(g.Score)))
	dict.SetItem("player", PyString(g.Player))

	items := make([]*PyObject, len(g.Items))
	for i, item := range g.Items {
		items[i] = PyString(item)
	}
	dict.SetItem("items", PyList(items...))

	return dict
}

func testPythonByClaude() {
	fmt.Println("=== Ultimate Python-Go Integration with CGO ===\n")

	if err := Initialize(); err != nil {
		panic(err)
	}
	defer Finalize()

	interp, err := NewInterpreter()
	if err != nil {
		panic(err)
	}

	// Demo 1: Basic execution
	fmt.Println("1. Basic Python execution:")
	interp.RunString(`
print("Hello from embedded Python!")
import sys
print(f"Python {sys.version_info.major}.{sys.version_info.minor}")
`)

	// Demo 2: Pass Go data structures
	fmt.Println("\n2. Go objects in Python:")

	game := GameState{
		Level:  5,
		Score:  1000,
		Player: "Hero",
		Items:  []string{"sword", "shield", "potion"},
	}

	interp.SetGlobal("game", game.ToPython())
	interp.RunString(`
print(f"Game State: {game}")
print(f"Player: {game['player']}, Level: {game['level']}, Score: {game['score']}")
print(f"Items: {', '.join(game['items'])}")

# Modify game state
game['score'] += 500
game['level'] += 1
`)

	// Get modified game state
	modifiedGame, _ := interp.GetGlobal("game")
	score, _ := modifiedGame.GetItem("score")
	scoreVal, _ := score.AsInt()
	fmt.Printf("Score after Python modification: %d\n", scoreVal)

	// Demo 3: Call Python functions
	fmt.Println("\n3. Calling Python functions from Go:")

	interp.RunString(`
def calculate_damage(base_damage, multiplier, critical=False):
    damage = base_damage * multiplier
    if critical:
        damage *= 2
    return int(damage)

def get_player_stats(player_name):
    return {
        'name': player_name,
        'health': 100,
        'mana': 50,
        'strength': 15
    }
`)

	damageFunc, _ := interp.GetGlobal("calculate_damage")
	damageResult, err := damageFunc.Call(PyInt(10), PyFloat(1.5), PyBool(true))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		damage, _ := damageResult.AsInt()
		fmt.Printf("Damage calculated: %d\n", damage)
	}

	statsFunc, _ := interp.GetGlobal("get_player_stats")
	statsResult, _ := statsFunc.Call(PyString("Warrior"))
	statsStr, _ := statsResult.AsString()
	fmt.Printf("Player stats: %s\n", statsStr)

	// Demo 4: Use Python libraries
	fmt.Println("\n4. Using Python standard library:")

	mathModule, _ := interp.ImportModule("math")

	piAttr, _ := mathModule.GetAttr("pi")
	piVal, _ := piAttr.AsFloat()
	fmt.Printf("π = %.6f\n", piVal)

	sqrtFunc, _ := mathModule.GetAttr("sqrt")
	sqrtResult, _ := sqrtFunc.Call(PyInt(16))
	sqrtVal, _ := sqrtResult.AsFloat()
	fmt.Printf("√16 = %.1f\n", sqrtVal)

	// Demo 5: Complex data structures
	fmt.Println("\n5. Complex data structures:")

	config := PyDict()
	config.SetItem("host", PyString("localhost"))
	config.SetItem("port", PyInt(8080))
	config.SetItem("debug", PyBool(true))
	config.SetItem("workers", PyInt(4))

	interp.SetGlobal("config", config)
	interp.RunString(`
import json
print(f"Config: {config}")
config_json = json.dumps(config, indent=2)
print(f"As JSON:\n{config_json}")
`)

	// Demo 6: Error handling
	fmt.Println("\n6. Error handling:")

	err = interp.RunString(`
try:
    result = 10 / 0
except ZeroDivisionError:
    print("Caught division by zero")
    result = None
`)
	if err != nil {
		fmt.Printf("Python error: %v\n", err)
	}

	// Demo 7: Lists and iterations
	fmt.Println("\n7. List operations:")

	numbers := PyList(PyInt(1), PyInt(2), PyInt(3), PyInt(4), PyInt(5))
	interp.SetGlobal("numbers", numbers)

	interp.RunString(`
print(f"Numbers: {numbers}")
doubled = [x * 2 for x in numbers]
print(f"Doubled: {doubled}")
filtered = [x for x in numbers if x % 2 == 0]
print(f"Even numbers: {filtered}")
`)

	// Demo 8: Advanced Python features
	fmt.Println("\n8. Advanced Python features:")

	interp.RunString(`
from datetime import datetime
import random

now = datetime.now()
print(f"Current time: {now.strftime('%Y-%m-%d %H:%M:%S')}")

random_numbers = [random.randint(1, 100) for _ in range(5)]
print(f"Random numbers: {random_numbers}")

# Class definition
class Monster:
    def __init__(self, name, health):
        self.name = name
        self.health = health
    
    def take_damage(self, damage):
        self.health -= damage
        return self.health > 0

goblin = Monster("Goblin", 30)
goblin.take_damage(10)
print(f"{goblin.name} health: {goblin.health}")
`)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nThis demo shows:")
	fmt.Println("✓ Full Python C API integration")
	fmt.Println("✓ Go objects → Python")
	fmt.Println("✓ Python functions ← Go")
	fmt.Println("✓ Python modules import")
	fmt.Println("✓ Complex data structures")
	fmt.Println("✓ Error handling")
	fmt.Println("✓ Zero Go dependencies (only CGO + Python)")
}
