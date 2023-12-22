use actix::prelude::*;


//Bez env loggera sa nevypisuje logy z async funkcii
use env_logger;
use std::fmt;
use std::{thread, time::Duration};

use super::warehouse::Warehouse;

pub struct App {
    pub time_step: u64,
    pub speed: u64,
    pub warehouse: Warehouse,
}

impl App {
    fn _new() -> App {
        App {
            warehouse: Warehouse::new(),
            time_step: 0,
            speed: 1,
        }
    }

    pub fn new() -> App {
        let mut app = App::_new();

        app.warehouse.set_resource("gold", 50);
        app.warehouse.set_resource("wood", 50);
        app.warehouse.set_resource("iron", 50);
        app.warehouse.set_resource("food", 50);
        app.warehouse.set_resource("stone", 50);

        return app;
    }

    pub fn step(&mut self) {
        println!("Doing cycle: {}", self.time_step);
        self.time_step += 1;
    }
}

impl Actor for App {

    type Context = Context<Self>;

    fn started(&mut self, ctx: &mut Context<Self>) {
        println!("Application started");

        loop {
            self.step();
            thread::sleep(Duration::from_secs(1));
        }
    }

    fn stopped(&mut self, ctx: &mut Context<Self>) {
        println!("Application is stopped");
    }
}

impl fmt::Display for App {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        writeln!(f, "{}", 0)
    }
}
