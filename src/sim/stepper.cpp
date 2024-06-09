#include "stepper.h"

#include "../prelude/prelude.h"


namespace simstudio {

    bool Stepper::IsEnd() {
        return false;
    }

    void Stepper::Step() {
        _active_time += _step_index;
        _step_index += 1;

        Log << "[Step][" << _active_time << "]" << Endl;
    }


}