import sys

START_CHAR = "["
END_CHAR = "]"


class Printer:
    passes = 1
    previous_page = -1
    previous_string = ""
    denominator = 0

    def emit(self, page_number: int) -> None:
        # set denominator if page_number < last_page_number
        if page_number < self.previous_page:
            self.denominator = self.previous_page
            self.passes += 1

        self.previous_page = page_number
        current_string = f"processing (pass={self.passes}): {page_number}"

        if self.denominator:
            current_string += f"/{self.denominator}"

        current_string = "\b" * len(self.previous_string) + current_string
        self.previous_string = current_string[len(self.previous_string) :]

        print(current_string, end="")
        sys.stdout.flush()

    def clear(self) -> None:
        print(
            "\b" * len(self.previous_string)
            + " " * len(self.previous_string)
            + "\b" * len(self.previous_string),
            end="",
        )
        sys.stdout.flush()


try:
    stored = ""
    printer = Printer()
    while True:
        buff = sys.stdin.read(1)

        # Exit once stdin is empty.
        if buff == "":
            printer.clear()
            break

        # Page numbers are delimited by square brackets like: [<number>].

        # Once the ] character is matched, emit the stored number
        if buff == END_CHAR and len(stored) > 1:
            try:
                printer.emit(int(stored[1:]))
            except ValueError:
                pass
            stored = ""

        # Save the number between square brackets
        elif len(stored) > 0 and buff != END_CHAR:
            if buff.isdigit():
                stored += buff
            else:
                stored = ""

        # Mark number start
        elif buff == START_CHAR:
            stored = START_CHAR

        # Empty stored buffer on mismatch
        else:
            stored = ""

except KeyboardInterrupt:
    exit(1)
