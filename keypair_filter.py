import json

# Path to your file
file_path = "keypairs.json"


def main():
    matching_key_ids = []

    # Open the file and read line by line (handles large files efficiently)
    with open(file_path, "r", encoding="utf-8") as f:
        # Load entire file as JSON array
        data = json.load(f)

        for entry in data:
            try:
                # Parse the serializedPublicKey JSON string
                pubkey = json.loads(entry.get("serializedPublicKey", "{}"))
                x_value = pubkey.get("x", "")

                # Check if length is less than 44
                if len(x_value) < 43:
                    matching_key_ids.append(entry.get("keyId"))

            except json.JSONDecodeError as e:
                print(f"Skipping entry due to JSON decode error: {e}")

    # Print all matching keyIds
    print("Matching keyIds:")
    matching_key_ids.sort()
    for kid in matching_key_ids:
        print(kid)


if __name__ == "__main__":
    main()
