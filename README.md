configure environment file
```bash
cp example_env .env
```

recompile tailwind css using the following command
```bash
npx tailwindcss -i ./ui/assets/source.css -o ./ui/assets/main.css --watch
```

start the project, make sure you have go air installed https://github.com/air-verse/air
if `air` command is not available on your system, make sure the go bin directory is in your path.

```bash
air
```
