use alloc::fmt::format;
use alloc::format;
use core::convert::Into;
use math::vec3::{Position};

use crate::inworld::InWorld;
use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;

use crate::instances_guard::{
    setup_instance_counter,
    get_and_increment_instance_count
};



setup_instance_counter!();


#[derive(Debug)]
pub struct Entity {
    pub name: String,
    pub uid: String,
    pub position: Position,
}


impl Entity {
    pub fn New() -> Self {

        Entity {
            name: "".into(),
            uid: format!("Entity_{}", get_and_increment_instance_count!()),
            position: Position::NewZero(),
        }
    }
}

impl InWorld for Entity {
    fn GetName(&self) -> String {
        self.name.clone()
    }

    fn GetUid(&self) -> String {
        self.uid.clone()
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

    fn Step(&mut self, stepper: &Stepper, context: &ToolsContext) {
        //Do something in step
    }

    fn Move(&mut self, position: &Position) {
        self.position.Copy(position)
    }
}