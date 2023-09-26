rm -f jmp
go get github.com/lithammer/fuzzysearch/fuzzy
go build
if [[ ! -f jmp ]]; then
  echo "Compilation failed.\n"
  exit 1
fi
echo "Compilation completed."

mkdir -p $HOME/.jmp/
mkdir -p $HOME/bin/

if [[ ! -f "$HOME/.jmp/cfg.json" ]]; then
  cp default_config.json $HOME/.jmp/cfg.json;
  echo "Copied default configuration to $HOME/.jmp/cfg.json";
fi

mv jmp $HOME/bin/jmpbin
echo "Moved jmp binary to $HOME/bin/jmpbin."
echo "Add the following to your \"~/.zshrc\":\n"

echo "
function jmp() {
  dir=\"\$(\$HOME/bin/jmpbin \$1)\"
  if [ ! -z \"\$dir\" ]; then
    chdir \"\$dir\";
  fi
}
"

echo "\nand then run \"source ~/.zshrc\" to get started."
