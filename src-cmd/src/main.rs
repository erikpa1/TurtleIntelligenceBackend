#![allow(warnings)]

extern crate alloc;

use alloc::boxed::Box;
use alloc::format;
use alloc::rc::Rc;
use core::any::Any;
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
        project.AddEntity(Rc::new(RefCell::new(Box::new(station))));
        // project.stations.push(Rc::new(RefCell::new(station)))
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

    project.Init(&stepper);

    while stepper.IsEnd() == false {
        project.Step(&stepper);
        stepper.Step();


        if stepper.step_index == 20 {
            println!("Adding entity to station");

            if let Some(station_ref) = project.entities_all.get_mut("Station_1".into()) {
                let tmp = station_ref.borrow_mut();

                println!("{:?}", tmp);



                // println!("Is my_type InWorld? {}", inworld_mut.is::<Station>());

                // println!("X - {}", inworld_mut.GetName());
                // if let Some(station) = inworld_mut.downcast_ref::<Station>() {
                //     // The entity is a Station
                //     println!("The entity is a Station: {:?}", station);
                // }
            }


            // project.stations[0].borrow_mut().TakeEntity(&project.entities[0]);
        }
    }


    project.PrintStatistics(&stepper);

    println!("Simulation ended");
}

// TODO list:
// spravit stav zasob, materialov
// spravit kriticke ukazovatele

