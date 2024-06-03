use crate::_std::*;


use crate::entity::Entity;
use crate::station::Station;
use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;
use crate::inworld::InWorld;

use std::any::Any;

pub struct Project {
    pub tools_context: ToolsContext,
    pub entities_all: HashMap<String, Mrc<dyn InWorld>>,
    pub entities: Vec<Rc<RefCell<Entity>>>,
    pub stations: Vec<Rc<RefCell<Station>>>,
}

impl Project {
    pub fn New() -> Self {
        Project {
            tools_context: ToolsContext::New(),
            entities: vec![],
            stations: vec![],
            entities_all: HashMap::new(),
        }
    }
    pub fn Step(&mut self, stepper: &Stepper) {
        for entity in &mut self.entities {
            entity.borrow_mut().Step(stepper, &mut self.tools_context);
        }
        for station in &mut self.stations {
            station.borrow_mut().Step(stepper, &mut self.tools_context);
        }
    }

    pub fn Init(&mut self, steper: &Stepper) {
        for entity in &mut self.entities {
            entity.borrow_mut().Init(steper, &self.tools_context)
        }
        for station in &mut self.stations {
            station.borrow_mut().Init(steper, &self.tools_context)
        }
    }

    pub fn PrintStatistics(&self, stepper: &Stepper) {
        for station in &self.stations {
            station.borrow().PrintStatistics(stepper)
        }
    }

    pub fn AddEntity(&mut self, entity: Mrc<dyn InWorld>) {
        let uid = entity.borrow().GetUid();
        self.entities_all.insert(uid, entity);
    }
}