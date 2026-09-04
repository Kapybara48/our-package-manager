use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command()]
struct Args {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Install { url: String },
    Remove,
    Update,
}

fn main() {
    let args = Args::parse();

    match args.command {
        Commands::Install { url } => println!("installing {}", url),
        Commands::Remove => println!("removing"),
        Commands::Update => println!("updating"),
    }
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
