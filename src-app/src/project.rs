use alloc::rc::Rc;
use alloc::vec::Vec;
use core::cell::RefCell;

use crate::entity::Entity;
use crate::station::Station;
use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;

pub struct Project {
    pub tools_context: ToolsContext,
    pub entities: Vec<Rc<RefCell<Entity>>>,
    pub stations: Vec<Rc<RefCell<Station>>>,
}

impl Project {
    pub fn New() -> Self {
        Project {
            tools_context: ToolsContext::New(),
            entities: vec![],
            stations: vec![],
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

    pub fn Init(&mut self) {
        for entity in &mut self.entities {
            entity.borrow_mut().Init()
        }
        for station in &mut self.stations {
            station.borrow_mut().Init()
        }
    }

    pub fn PrintStatistics(&self, stepper: &Stepper) {
        for station in &self.stations {
            station.borrow().PrintStatistics(stepper)
        }
    }
}