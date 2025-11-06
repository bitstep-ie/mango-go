# mango4go - random

The `random` package provides random helpers


## Number 
Returns a random number in the range [min, max] inclusive for int and [min, max) for floats. Safely swaps max & min
### Limitations
  - uses math/rand (not safe for concurrent use, not for security-sensitive use)
  - overflow warning on extreme ranges


## Sign
As simple as returning either +1 or -1 randomly


## Bool 
Randomly returns either true or false


## Choice
Picks a random element from a non-empty slice. **Panics** if the slice is `empty`.


## Byte
Returns a single secure random byte [0,255]


## String 
You get a random alphanumeric string of length n. Includes both `lowercase` and `uppercase` letters.


## Alpha 
Get random letters - length n. Includes both `lowercase` and `uppercase` letters.


## Numeric
Get a random string of digits of length n.


## FromCharset
Returns a random string from your specific charset


## Password 
Generates a random password of length n according to provided options. 
**Panics** if no charsets selected (or effectively not selected - e.g: excluding all options)

### PasswordOptions
  - Letters: a-zA-Z
  - Digits: 0-9
  - Symbols: !@#$%^&*()-_=+[]{}<>?/|~
  - Exclude: all strings you don't want to take part in password generation


## Date 
Random `time.Time` between `[min, max]` (inclusive)
*Note:* Safely swaps `min` and `max` if `min>max`.


## Duration 
Random `time.Duration` between `[min, max]` (inclusive)
*Note:* Safely swaps `min` and `max` if `min>max`.

