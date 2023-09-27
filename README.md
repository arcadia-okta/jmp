# jmp

A simple command-line utility that will `chdir` into your favourite git repositories, using fuzzy
matching to save you time typing.

## Installation

The `install.sh` script will take care of installing dependencies and building `jmp` as well as
moving everything into the right places.  To really get going, though, you'll want to define a
function in your `$HOME/.zshrc` or equivalent to make it work properly.

```bash
function jmp() {
  dir="$($HOME/bin/jmpbin $1)"
  if [ ! -z "$dir" ]; then
    chdir "$dir";
  fi
}
```

Once you've added the above, `source $HOME/.zshrc` and you're good to go!

## Configuration

After being installed, a configuration file will live in `$HOME/.jmp/cfg.json` on your filesystem.
Here you can specify the root directory for `jmp` to search from.  If all of your cloned
repositories live under a certain directory, you should modify the `projectsRoot` value
appropriately.

`jmp` uses fuzzy matching to find directories. Due to the nature of these types of algorithms,
you will sometimes not get what you want and find yourself in a completely unexpected directory.
To help remedy this, it allows you to declare your favourite repositories and will bias heavily
in favour of them.

The `searchDepth` configuration option helps to keep searching fast.  The lower this value is,
the faster `jmp` will work. The higher it is, the deeper into nested directories it will search.

## Usage

Once you're happy with your configuration, just start `jmp`ing around!

```
jmp <fuzzy-string>
```
