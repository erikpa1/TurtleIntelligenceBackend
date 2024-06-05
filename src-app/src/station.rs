use alloc::rc::Rc;
use alloc::{format, vec};
use alloc::fmt::format;
use alloc::string::String;
use alloc::vec::Vec;
use core::convert::Into;
use core::cell::RefCell;
use core::clone::Clone;
use math::time::TimeExpresionExecutioner;

use math::vec3::Position;
use crate::entity::Entity;
use crate::inworld::InWorld;

use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;

use std::sync::atomic::{AtomicUsize, Ordering};
use lazy_static::lazy_static;



#[derive(Debug, PartialEq, Eq)]
pub enum NotEnoughtEntitiesMode {
    NON_OPERATIVE,
    HANDICAPED,
}

#[derive(Debug, PartialEq, Eq)]
pub enum VariableType {
    CONST,
    FUNCTION,
}

lazy_static! {
    static ref INSTANCE_COUNT: AtomicUsize = AtomicUsize::new(0);
}

#[derive(Debug)]
pub struct Station {
    pub name: String,
    pub uid: String,
    pub required_entities: u64,
    pub not_enought_entities_mode: NotEnoughtEntitiesMode,
    pub operation_time: String,
    pub operation_time_vt: VariableType,
    pub handicap_function: String,
    pub position: Position,
    pub entities: Vec<Rc<RefCell<Entity>>>,
    pub manufactured_count: u64,
    pub manufacturing_start: f64,
    pub is_manufacturing: bool,
    pub working_time: f64,
}

impl Station {
    pub fn GetInstancesCount() -> usize {
        INSTANCE_COUNT.load(Ordering::SeqCst)
    }

    pub fn New() -> Station {
        INSTANCE_COUNT.fetch_add(1, Ordering::SeqCst);

        Station {
            name: "".into(),
            uid: format!("Station_{}", Self::GetInstancesCount()),
            required_entities: 0,
            not_enought_entities_mode: NotEnoughtEntitiesMode::HANDICAPED,
            handicap_function: "".into(),
            operation_time: "00:00".into(),
            operation_time_vt: VariableType::CONST,
            position: Position::NewZero(),
            entities: vec![],
            manufactured_count: 0,
            manufacturing_start: 0.0,
            is_manufacturing: false,
            working_time: 0.0,
        }
    }

    pub fn TakeEntity(&mut self, entity: &Rc<RefCell<Entity>>) {
        entity.borrow_mut().position.Copy(&self.position);
        self.entities.push(entity.clone());
    }


    pub fn Step(&mut self, stepper: &Stepper, context: &mut ToolsContext) {
        // println!("Entity doing something random [{}]", context.expr.Execute(&"standard()".into()));
        self._0_CheckWorkers(stepper);
    }

    fn _0_CheckWorkers(&mut self, stepper: &Stepper) {
        if self.required_entities == 0 {
            self._1_CheckManufacturing(stepper);
        } else {
            if self.entities.len() == 0 {
                return;
            } else {
                if self.not_enought_entities_mode == NotEnoughtEntitiesMode::NON_OPERATIVE {
                    return;
                } else if self.not_enought_entities_mode == NotEnoughtEntitiesMode::HANDICAPED {
                    self._1_CheckManufacturing(stepper);
                }
            }
        }
    }

    fn _1_CheckManufacturing(&mut self, stepper: &Stepper) {
        if self.is_manufacturing == false {
            self._1_StartManufacturing(stepper);
        } else {
            println!("Start: [{}]", self.manufacturing_start);

            let man_duration = self.manufacturing_start + TimeExpresionExecutioner::SecondsFromTimeString(&self.operation_time) as f64;

            if man_duration < stepper.active_time {
                self.is_manufacturing = false;
                self.manufactured_count += 1;
                self._1_StartManufacturing(stepper);
            }
        }
    }

    fn _1_StartManufacturing(&mut self, stepper: &Stepper) {
        self.manufacturing_start = stepper.active_time;
        self.is_manufacturing = true;
    }


    pub fn PrintStatistics(&self, stepper: &Stepper) {
        println!("Station [{}]:", self.name);
        println!("Manufactured [{}]:", self.manufactured_count);
        println!("Occupied [{}%]:", (self.working_time / stepper.finish_time) * 100.0);
    }
}


impl InWorld for Station {
    fn GetName(&self) -> String {
        return self.name.clone();
    }

    fn GetUid(&self) -> String {
        return self.uid.clone()
    }

    fn GetType(&self) -> String {
        "entity".into()
    }

    fn Init(&mut self, step: &Stepper, context: &ToolsContext) {
        println!(
            "File: {}, Line: {}, {}",
            file!(),
            line!(),
            format!("[{}] received init", &self.name)
        );
    }
}


