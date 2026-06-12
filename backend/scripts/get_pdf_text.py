import sys
import json
import fitz

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No file path provided"}))
        sys.exit(1)

    file_path = sys.argv[1]

    try:
        doc = fitz.open(file_path)
        full_text = ""
        for page in doc:
            full_text += page.get_text()
        doc.close()

        if not full_text.strip():
            print(json.dumps({"error": "No text found in PDF. It may be a scanned image."}))
            return

        print(json.dumps({"text": full_text}))

    except Exception as e:
        print(json.dumps({"error": "Could not read PDF file."}))

main()