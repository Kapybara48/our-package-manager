use clap::{Parser, Subcommand};

mod build;
mod error;
mod git;
mod paths;

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
        Commands::Install { url } => install(&url),
        Commands::Remove => println!("removing"),
        Commands::Update => println!("updating"),
    }
}

fn install(url: &str) {
    println!("installing {}", url);
    match git::clone(url) {
        Ok(exit_status) => {
            if exit_status.success() {
                println!("successfully cloned");
            } else {
                println!("cloning was not successfull");
            }
        }
        Err(error) => println!("{}", error),
    }
}
