//
// Created by erikp on 5. 6. 2024.
//

#include <iostream>

#include "src/prelude/prelude.h"

#include "src/test.h"
#include "src/app/app.h"
#include "src/sim/stepper.h"
#include "src/sim/entity.h"

using namespace simstudio;

int main() {


    App app;

    Stepper stepper;
    stepper._finish_time = 100;


    while (stepper.IsEnd() == false) {
        stepper.Step();
    }

    Log << "Hello world " << TESTVAR << Endl;
    Log << "Should be red " << TESTVAR << Endl;

    return 0;
}