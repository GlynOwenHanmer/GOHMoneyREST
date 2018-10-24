### Disclaimer:
This repository contains some not-so-best-practices as it is a personal project in which I try out new techniques and ways of working. As always, some of these techniques end up not being the best way of working; they are kept in the project, however, until I find a good time/requirement to refactor.

# mon

Mon is a simple finance-tracking set of packages and application.

### Future Improvements
- Remove package globals
- Implement reusable tables
	- I have a design prepared for this using first-class functions and closures, which I believe could be pretty good. Hopefully the developer experience doesn't end up being too complex.
- Implement Go modules
- Move to an event-sourced model
- Improve error handling of `moncli`
	- It would be best to have error feedback presented in a more user-friendly way for `moncli` 