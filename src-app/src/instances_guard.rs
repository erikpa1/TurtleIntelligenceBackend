macro_rules! setup_instance_counter {
    () => {
        use std::sync::atomic::{AtomicUsize, Ordering};
        use lazy_static::lazy_static;

        lazy_static! {
            static ref INSTANCE_COUNT: AtomicUsize = AtomicUsize::new(0);
        }
    };
}
pub(crate) use setup_instance_counter;


macro_rules! increment_instance_count {
    () => {
        INSTANCE_COUNT.fetch_add(1, Ordering::SeqCst);
    };
}
pub(crate) use increment_instance_count;


macro_rules! get_instance_count {
    () => {
        INSTANCE_COUNT.load(Ordering::SeqCst)
    };
}

pub(crate) use get_instance_count;


macro_rules! get_and_increment_instance_count {
    () => {
        { let index = INSTANCE_COUNT.load(Ordering::SeqCst);
        INSTANCE_COUNT.fetch_add(1, Ordering::SeqCst);
        index
        }
    };
}

pub(crate) use get_and_increment_instance_count;
