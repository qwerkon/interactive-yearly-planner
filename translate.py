from typing import Callable
import json
from sys import argv
from glob import glob
from os.path import exists

language = argv[1].lower()
TRANSLATION_FOLDER = "translations/"
file = f"{TRANSLATION_FOLDER}{language}.json"

if file in glob(f"{TRANSLATION_FOLDER}*.json"):
    with open(file, "r") as f:
        translation = json.load(f)
else:
    raise ValueError("Requested translation is not currently supported.\nThe program will now exit.")


if any(not word.isascii() for word in translation.values()):
    print("unicode found")
    font_edit = r""
else:
    font_edit = ""
print(f"Translating pdf to {language}")

MONTHS = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
MONTHS_SHORT = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
WEEKDAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
WEEKDAYS_SHORT = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
WEEKDAYS_SINGLE_HEADER = ("W & M & T & W & T & F & S & S", "T & P & W & Ś & C & P & S & N")
WEEK = ["Week"]
CALENDAR = ["Calendar"]
NOTES = ["Notes"]
NOTE = ["Note"]
NOTES_INDEX = ["Notes Index"]
ALL_NOTES = ["All notes"]
SCHEDULE = ["Schedule"]
PRIORITIES = ["Top priorities"]
MORE = ["More"]
REFLECT = ["Reflect"]
PHRASES = ["Things I'm grateful for", "The best thing that happened today", "Daily log"]

def handle_all() -> None:
    if exists("out/annual.tex"):
        handle_annual()
    if exists("out/quarterly.tex"):
        handle_quarterly()
    if exists("out/monthly.tex"):
        handle_monthly()
    if exists("out/weekly.tex"):
        handle_weekly()
    if exists("out/daily.tex"):
        handle_daily()
    if exists("out/daily_reflect.tex"):
        handle_daily_reflect()
    if exists("out/daily_notes.tex"):
        handle_daily_notes()
    if exists("out/notes_indexed.tex"):
        handle_notes_indexed()
    handle_common_generated_text()

def add_identifier(keys: list[str], func: Callable[[str], str] = lambda x: x, dictionary: dict[str, str] = translation) -> dict[str, str]:
    return {func(key): func(dictionary.get(key)) for key in keys}

def handle_annual() -> None:
    with open("out/annual.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "{" + x + "}}")
    replace |= add_identifier(MONTHS, lambda x: "}{" + x + "}")
    replace |= add_identifier(NOTES, lambda x: "{" + x + "}")
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/annual.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_quarterly() -> None:
    with open("out/quarterly.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "{" + x + "}}")
    replace |= add_identifier(NOTES, lambda x: "{" + x + "}")
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/quarterly.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_monthly() -> None:
    with open("out/monthly.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "}{" + x + "}")
    replace |= add_identifier(WEEKDAYS)
    replace |= add_identifier(WEEK, lambda x: "[c]{" + x)
    replace |= add_identifier(NOTES, lambda x: "{" + x + "}")
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/monthly.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_weekly() -> None:
    with open("out/weekly.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "}{" + x + "}")
    replace |= add_identifier(WEEK, lambda x: "}{" + x)
    replace |= add_identifier(WEEKDAYS, lambda x: ", " + x + "}")
    replace |= add_identifier(NOTES, lambda x: "{" + x)
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/weekly.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_daily() -> None:
    with open("out/daily.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "}{" + x + "}")
    replace |= add_identifier(WEEK, lambda x: "}{" + x)
    replace |= add_identifier(WEEKDAYS, lambda x: "}{" + x + ",")
    replace |= add_identifier(WEEKDAYS_SHORT, lambda x: "}{" + x + ",")
    replace |= add_identifier(SCHEDULE, lambda x: "{" + x + "\\")
    replace |= add_identifier(PRIORITIES, lambda x: "{" + x + "\\")
    replace |= add_identifier(NOTES, lambda x: "{" + x + " $")
    replace |= add_identifier(MORE, lambda x: "{" + x + "}")
    replace |= add_identifier(REFLECT, lambda x: "{" + x + "}")
    replace |= add_identifier(ALL_NOTES)
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/daily.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_daily_reflect() -> None:
    with open("out/daily_reflect.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "}{" + x + "}")
    replace |= add_identifier(WEEK, lambda x: "}{" + x)
    replace |= add_identifier(WEEKDAYS, lambda x: "}{" + x + ",")
    replace |= add_identifier(WEEKDAYS_SHORT, lambda x: "}{" + x + ",")
    replace |= add_identifier(REFLECT, lambda x: "{" + x + "}")
    replace |= add_identifier(PHRASES)
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/daily_reflect.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_daily_notes() -> None:
    with open("out/daily_notes.tex", "r") as file:
        text = file.read()

    replace = add_identifier(MONTHS, lambda x: "}{" + x + "}")
    replace |= add_identifier(WEEK, lambda x: "}{" + x)
    replace |= add_identifier(WEEKDAYS, lambda x: "}{" + x + ",")
    replace |= add_identifier(WEEKDAYS_SHORT, lambda x: "}{" + x + ",")
    replace |= add_identifier(NOTES, lambda x: "{" + x + "}")
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/daily_notes.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_notes_indexed() -> None:
    with open("out/notes_indexed.tex", "r") as file:
        text = file.read()

    replace = add_identifier(NOTES_INDEX, lambda x: "}{" + x)
    replace |= add_identifier(NOTE, lambda x: "}{" + x)
    for english, spanish in replace.items():
        text = text.replace(english, spanish)

    with open("out/notes_indexed.tex", "w") as file:
        file.write(font_edit)
        file.write(text)

def handle_common_generated_text() -> None:
    """Translate labels generated by Go after templates are rendered."""
    paths = glob("out/*.tex")

    for path in paths:
        with open(path, "r") as file:
            text = file.read()

        replace = {}
        replace |= add_identifier(CALENDAR, lambda x: "}{" + x + "}")
        replace |= add_identifier(MONTHS_SHORT, lambda x: "}{" + x + "}")
        replace |= add_identifier(NOTES, lambda x: "}{" + x + "}")
        replace |= add_identifier(MONTHS, lambda x: " & " + x + " &")
        replace |= add_identifier(WEEKDAYS, lambda x: r"\textbf{" + x + "}")
        replace |= add_identifier(WEEKDAYS, lambda x: " & " + x + " &")
        replace |= add_identifier(WEEKDAYS_SHORT, lambda x: "{" + x + ",")
        replace |= add_identifier(WEEKDAYS_SHORT, lambda x: "}{" + x + ",")

        for english, translated in replace.items():
            text = text.replace(english, translated)

        text = text.replace(*WEEKDAYS_SINGLE_HEADER)
        text = text.replace("Notatkas", translation.get("Notes", "Notes"))
        text = text.replace("{Notatki Index}", "{Notes Index}")
        text = text.replace("{Indeks notatek}", "{Notes Index}")

        with open(path, "w") as file:
            file.write(text)

handle_all()
