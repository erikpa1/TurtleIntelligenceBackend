use core::clone::Clone;

pub struct Stepper {
    pub active_time: f64,
    pub finish_time: f64,
    pub step: f64,
    pub last_offset: f64,
    pub step_index: u64
}

impl Stepper {
    pub fn New() -> Self {
        Stepper {
            active_time: 0.0,
            finish_time: 0.0,
            last_offset: 0.0,
            step: 1.0,
            step_index: 0
        }
    }

    pub fn Step(&mut self) {
        let previous = self.active_time.clone();
        self.active_time += self.step;
        self.last_offset = self.active_time - previous;
        self.step_index += 1;
    }

    pub fn IsEnd(&self) -> bool {
        if (self.finish_time == 0.0) {
            return false;
        } else {
            self.active_time > self.finish_time
        }
    }
}