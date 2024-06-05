macro_rules! to_string {
    ($self:expr) => {
        format!("[{}, {}, {}]", $self.x, $self.y, $self.z)
    };
}


#[derive(Debug)]
pub struct Vec3<T> {
    pub x: T,
    pub y: T,
    pub z: T,
}


impl<T> Vec3<T> where
    T: Clone,
    T: std::fmt::Display,
{
    pub fn New(x: T, y: T, z: T) -> Self {
        return Self {
            x,
            y,
            z,
        };
    }

    pub fn Copy(&mut self, another: &Vec3<T>) {
        self.x = another.x.clone();
        self.y = another.y.clone();
        self.z = another.z.clone();
    }


    pub fn ToString(&self) -> String {
        format!("[{}, {}, {}]", self.x, self.y, self.z)
        // return to_string!(self);
    }
}

impl Vec3<f32> {
    pub fn NewZero() -> Self {
        Vec3 {
            x: 0.0,
            y: 0.0,
            z: 0.0,
        }
    }
}

pub type Position = Vec3<f32>;
pub type Scale = Vec3<f32>;