

namespace simstudio {

    class Stepper {

    public:
        double _active_time = 0;
        double _finish_time = 0;
        double _last_offset = 0;
        double _step = 1;
        double _step_index = 0;


        void Step();
        bool IsEnd();

    };



}