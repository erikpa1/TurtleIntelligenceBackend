mod vec3;

use vec3::*;

fn main() {
    let mut tmp = Position::NewZero();

    println!("Zero vector: {}", tmp.ToString());

    let newOne = Position::New(5.0, 5.0, 5.0);
    tmp.Copy(&newOne);

    println!("After copy: {}", tmp);
}