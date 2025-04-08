use core::convert::Into;
use math::vec3::Position;
use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;

use std::any::{Any, TypeId};

use core::fmt::Debug;

pub trait InWorld: Any + Debug {


    fn GetPosition(&self) -> Position {
        return Position::NewZero();
    }

    fn Move(&mut self, position: &Position) {}

    fn GetType(&self) -> String {
        return "".into();
    }

    fn GetName(&self) -> String {
        return "".into();
    }

    fn GetUid(&self) -> String {
        return "".into();
    }

    fn Step(&mut self, stepper: &Stepper, context: &ToolsContext) {
        //override
    }

    fn Init(&mut self, step: &Stepper, context: &ToolsContext) {
        //override
    }
}

