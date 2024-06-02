extern crate alloc;

use alloc::format;
use alloc::rc::Rc;
use core::cell::RefCell;
use core::convert::Into;
use math::*;
use app::*;

use app::entity::Entity;
use app::project::Project;
use app::station::NotEnoughtEntitiesMode::HANDICAPED;
use app::station::Station;
use app::stepper::Stepper;
use math::expr::MathExpresionExecutioner;
use math::time::TimeExpresionExecutioner;


fn main() {
    let mut stepper = Stepper::New();
    stepper.finish_time = 100.0;

    let mut project = Project::New();

    for index in 0..3 {
        let mut station = Station::New();
        station.name = format!("Station_{}", index);

        let mut rnd = MathExpresionExecutioner::New();

        let milis = rnd.Execute(&"standard(5000, 10000)".into()) as u64;
        let operation_time = TimeExpresionExecutioner::MakeFromMilis(milis);

        station.operation_time = operation_time.into();
        project.stations.push(Rc::new(RefCell::new(station)))
    }

    // let station0 = project.stations[0].borrow_mut();
    // station0.required_entities = 2;
    // station0.not_enought_entities_mode = HANDICAPED;
    // station0.operation_time = "00:05".into();


    for index in 0..5 {
        let mut entity = Entity::New();
        entity.name = format!("Entity_{}", index);
        project.entities.push(Rc::new(RefCell::new(entity)));
    }

    println!("Count: {}", project.entities.len());

    project.Init();

    while stepper.IsEnd() == false {
        project.Step(&stepper);
        stepper.Step();


        if stepper.step_index == 20 {
            println!("Adding entity to station");
            project.stations[0].borrow_mut().TakeEntity(&project.entities[0]);
        }

        if stepper.step_index == 30 {
            println!("-----------");
            println!("Machine entities: ");

            let station = project.stations[0].borrow();

            for entity in &station.entities {
                let mut entity = entity.borrow();
                println!("Entity: {}", entity.name);
            }
            println!("-----------");
        }
    }


    project.PrintStatistics(&stepper);

    println!("Simulation ended");
}

// TODO list:
// spravit stav zasob, materialov
// spravit kriticke ukazovatele

