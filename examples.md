# Shamir CLI Usage Examples

## Basic Commands

### Secret Splitting

```bash
# Split a secret into 5 parts, 3 parts needed for recovery
./shamir-cli split "My secret password" 5 3

# Result:
# Secret split into 5 parts, 3 parts required for recovery:
# Part 1: 1:a98aecae01fe585fb3b74ff8256028697f84681355b932fadc7f5260ec19ee3dc681af
# Part 2: 2:8239560cf29a3978a7e63010930ceec867e16495e64cc4a765446fe7135624010f0f77
# ... other parts
```

### Secret Recovery

```bash
# Recover secret from any 3 parts
./shamir-cli combine "1:a98aecae01fe585fb3b74ff8256028697f84681355b932fadc7f5260ec19ee3dc681af,3:fb126a1723deb0a7c4e4ae6a66d1172ac8dc2c560c25468c39eb8357449e461cf8bceb,5:bbb42ecd1f283adbaf56a3445d5923a9efbc9dd489b96d966225a325bf3c0e13bbe0a0"

# Result:
# Recovered secret: My secret password
```

## Practical Examples

### Example 1: Password Backup

```bash
# Create password backup with 7 parts, 4 needed for recovery
./shamir-cli split "SuperSecretPassword123!" 7 4

# Result allows distributing parts to 7 different people
# Any 4 of them can recover the password
```

### Example 2: Cryptographic Key Protection

```bash
# Split key into 10 parts, 6 needed for recovery
./shamir-cli split "abcdef123456789" 10 6

# Parts can be stored in different locations for maximum security
```

### Example 3: Minimal Scheme

```bash
# Simplest scheme: 3 parts, 2 needed
./shamir-cli split "test" 3 2

# Result:
# Part 1: 1:4b8e5927
# Part 2: 2:63a71c45
# Part 3: 3:7bc56f23

# Recovery with any two parts:
./shamir-cli combine "1:4b8e5927,3:7bc56f23"
# Result: test
```

## Important Features

### Security
- Any number of parts below threshold cannot recover the secret
- Parts look like random data
- Each run creates different parts for the same secret

### Limitations
- Maximum 255 parts
- Minimum 2 parts for recovery
- Total number of parts must be >= threshold

### Testing
```bash
# Built-in algorithm test
./shamir-cli test

# Result will show algorithm success
```

## Usage Scenarios

1. **Corporate Security**: Split administrator password between multiple employees
2. **Personal Security**: Backup important passwords
3. **Cryptography**: Protect private keys
4. **Family Security**: Access to important data only with participation of multiple family members

## Troubleshooting

### Invalid Part Format
```bash
# Wrong:
./shamir-cli combine "1-abcdef,2-123456"

# Correct:
./shamir-cli combine "1:abcdef,2:123456"
```

### Insufficient Parts
```bash
# If 3 parts needed but only 2 provided, result will be incorrect
./shamir-cli combine "1:abcdef,2:123456"  # Incorrect result
```