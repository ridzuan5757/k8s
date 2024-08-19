import random
import string

# Generate a 24-character random master key
master_key = ''.join(random.choices(
    string.ascii_letters + string.digits, k=24))

# Print the master key
print("Generated master key:", master_key)
