pub struct Industry {
    pub farms: Vec<Farm>,
}

impl Industry {
    pub fn new() -> Industry {
        Industry { farms: Vec::new() }
    }

    pub fn add_farm(&mut self) {
        self.farms.push(Farm::new())
    }
}

pub struct Farm {
    pub efectivity: f32,
}

impl Farm {
    pub fn new() -> Farm {
        Farm { efectivity: 2.5 }
    }
}
