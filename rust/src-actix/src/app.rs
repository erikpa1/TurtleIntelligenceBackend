use actix::prelude::*;
use serde_json::{json, Value};

use std::borrow::BorrowMut;

use std::fmt;
use std::sync::{Arc, Mutex};
use std::{thread, time::Duration};

use super::constans;
use super::industry::Industry;
use super::warehouse::Warehouse;

pub struct App {
    pub time_step: u64,
    pub speed: u64,
    pub population: u64,
    pub warehouse: Warehouse,
    pub industry: Industry,

}

impl App {
    fn _new() -> App {
        App {
            time_step: 0,
            speed: 1,
            population: 4,

            warehouse: Warehouse::new(),
            industry: Industry::new(),
        }
    }

    pub fn new() -> App {
        let mut app = App::_new();

        app.warehouse.set_resource(constans::GOLD, 50.0);
        app.warehouse.set_resource(constans::WOOD_LOGS, 50.0);
        app.warehouse.set_resource(constans::IRON, 50.0);
        app.warehouse.set_resource(constans::FOOD, 50.0);
        app.warehouse.set_resource(constans::STONE, 50.0);

        app.industry.add_farm();
        app.industry.add_farm();

        return app;
    }

    pub fn step(&mut self) {
        println!("Doing cycle: {}", self.time_step);
        self.time_step += 1;

        for x in &self.industry.farms {
            self.warehouse.add_resouce("food", x.efectivity)
        }

        self.warehouse
            .lower_resource("food", 1.0 * self.population as f32)
    }

    pub fn to_json(&self) -> Value {
        json!({
            "time_step": self.time_step,
            "speed": 1,
            "warehouse": self.warehouse.to_json()
        })
    }
}

pub struct AppActor {
    app: Arc<Mutex<App>>,
}

impl AppActor {
    pub fn new(app: &Arc<Mutex<App>>) -> AppActor {
        AppActor { app: app.clone() }
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

pub type AppArcMutext = Arc<Mutex<App>>;
