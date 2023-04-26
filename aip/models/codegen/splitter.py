import pathlib
import re
import shutil

BEGIN_FILE_MARKER = re.compile(r'OUTPUT FILE +([a-zA-Z0-9_./ -]+):')
END_FILE_MARKER = re.compile(r'LLM EOF')
EOM_MARKER = re.compile(r'END OF MESSAGE')

class FileSplitter():
    def __init__(self, output_path, immediate=False, clean=False):
        self.immediate = immediate
        self.current_file = ""
        self.output_path = pathlib.Path(output_path)
        self.files = {}

        if clean:
            shutil.rmtree(str(self.output_path), ignore_errors=True)

        self.output_path.mkdir(parents=True, exist_ok=True)

    def parse(self, line):
        print(line)

        m = re.search(BEGIN_FILE_MARKER, line)

        if m is not None:
            self.begin_file(m.group(1))
            return

        m = re.search(END_FILE_MARKER, line)

        if m is not None:
            self.end_file()
            return

        self.append_line(line)

    def begin_file(self, name):
        if self.current_file != name:
            self.end_file()

        self.current_file = name

    def end_file(self):
        if self.current_file == "":
            return

        if self.immediate:
            self.emit_file(self.current_file)

        self.current_file = ""

    def append_line(self, line):
        if self.current_file != "":
            if self.current_file not in self.files:
                self.files[self.current_file] = ""

            self.files[self.current_file] += line + "\n"

    def emit_file(self, name):
        if not name in self.files:
            return

        contents = self.files[name]
        path = self.output_path.joinpath(name)

        path.parent.mkdir(parents=True, exist_ok=True)

        with path.open(mode="w") as f:
            f.write(contents)

    def emit(self):
        self.end_file()

        if not self.immediate:
            for file in self.files.keys():
                self.emit_file(file)
