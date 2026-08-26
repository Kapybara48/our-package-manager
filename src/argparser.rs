use std::env;

pub fn print_args() {
    let args = env::args().skip(1);

    for arg in args {
        println!("{}", arg);
    }
}
