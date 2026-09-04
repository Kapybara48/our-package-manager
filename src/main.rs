use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command()]
struct Args {
    #[arg(short)]
    name: String,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Install,
    Remove,
    Update,
}

fn main() {
    let args = Args::parse();

    println!("{}", args.name);
}

fn main_copy() {
    let mut args = std::env::args();
    match args.nth(1) {
        Some(arg) => subcommand(arg.as_str(), &mut args),
        None => println!("no option selected/help menu"),
    }
}
fn subcommand(arg: &str, args: &mut std::env::Args) {
    match arg {
        "install" => install(args),
        "remove" => println!("removing..."),
        "update" => println!("updating..."),
        _ => println!("unknown command"),
    };
}
fn install(args: &mut std::env::Args) {
    match args.next() {
        Some(arg) => println!("Installing {}", arg),
        None => println!("no package selected"),
    }
}
