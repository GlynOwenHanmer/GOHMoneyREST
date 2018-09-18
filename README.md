Making changes to this file so that I can test the security of the CI

### Disclaimer:
This repository contains some not-so-best-practices as it is a personal project in which I try out new techniques and ways of working. As always, some of these techniques end up not being the best way of working; they are kept in the project, however, until I find a good time to refactor.  

It's all part of the learning process, right?!


# mon

Mon is a simple finance-tracking set of packages and application.

### BUGS
- When updating an account, possibly also with other commands, when the opened and closed date are set to the same date, the command reports that the closed date is before the opened date
- Due to the flakey way that dates are stored, it's possible for issues to arise when having a closing date that is the same as a balance that exists. For example, if a last balance is inserted at 13h45 on 2018-05-04, trying to add a close date to the account of 2018-05-04 may cause an error if the date trying to be inserted has a time of 00h00. To sort this out, the dates and times of this whole thing will need to be sorted out. Perhaps it would just be best to use Date and not get involved with time at this stage.
