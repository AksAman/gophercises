# Exercise: Phone Number Normalizer

Original link : [gophercises/phone](https://github.com/gophercises/phone)

## Exercise details

### Problem statement
- Writing a program that 
    - will iterate through a database
    - normalize all of the phone numbers in the DB
    - Remove any duplicates

- Example input

    Phone numbers with 10 digits with or without dashes, spaces, or parentheses


    ```
    1234567890
    123 456 7891
    (123) 456 7892
    (123) 456-7893
    123-456-7894
    123-456-7890
    1234567892
    (123)456-7892
    ```

- Program should make the numbers match the format
    ```
    ##########
    ```
    That is, we are going to remove all formatting and only store the digits. When we want to display numbers later we can always format them, but for now we only need the digits.

- Output

    ```
    1234567890
    1234567891
    1234567892
    1234567893
    1234567894
    ---- was a duplicate, removed ----
    ---- was a duplicate, removed ----
    ---- was a duplicate, removed ----
    ```

### Learning Outcomes
- [x] Some string manipulation
- [x] Learn how to write raw SQL using the [`database/sql`](https://golang.org/pkg/database/sql/) package in the standard library
- [ ] Learn how to use the [`sqlx`](https://github.com/jmoiron/sqlx) package, which is an extension/wrapper around `database/sql` that makes it easier to work with.
- [ ] Learn how to use a minimalistic ORM. Here we'll use the [`gorm`](https://github.com/jinzhu/gorm) package
