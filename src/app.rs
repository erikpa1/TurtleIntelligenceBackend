use actix::prelude::*;


use std::borrow::BorrowMut;
use std::sync::Mutex;
//Bez env loggera sa nevypisuje logy z async funkcii
use env_logger;
use std::fmt;
use std::sync::Arc;
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
        // println!("Doing cycle: {}", self.time_step);
        self.time_step += 1;
        self.warehouse.lower_resource("gold", 1)
    }
}

pub struct AppActor {
    app: Arc<Mutex<App>>,
}

impl AppActor {
    pub fn new(app: &Arc<Mutex<App>>) -> AppActor {
        AppActor {
            app: app.clone()
        }
    }

    pub fn step(&mut self) {

        let mut appResult = self.app.lock();

        if let Ok(mut app) = appResult {
            app.step();
        }
        

    }
}

impl Actor for AppActor {
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
