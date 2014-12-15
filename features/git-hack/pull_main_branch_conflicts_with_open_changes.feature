Feature: git-hack handling conflicting remote main branch updates with open changes

  Background:
    Given I have a feature branch named "existing_feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | remote   | conflicting remote commit | conflicting_file | remote content |
      |        | local    | conflicting local commit  | conflicting_file | local content  |
    And I am on the "existing_feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack new_feature` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH           | COMMAND                |
      | existing_feature | git stash -u           |
      | existing_feature | git checkout main      |
      | main             | git fetch --prune      |
      | main             | git rebase origin/main |
    And my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git hack --abort`
    Then it runs the Git commands
      | BRANCH           | COMMAND                       |
      | HEAD             | git rebase --abort            |
      | main             | git checkout existing_feature |
      | existing_feature | git stash pop                 |
    And I end up on the "existing_feature" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And I have the following commits
      | BRANCH | LOCATION | MESSAGE                   | FILES            |
      | main   | remote   | conflicting remote commit | conflicting_file |
      |        | local    | conflicting local commit  | conflicting_file |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git hack --continue` while allowing errors
    Then I get the error "You must resolve the conflicts before continuing the git hack"
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git hack --continue `
    Then it runs the Git commands
      | BRANCH      | COMMAND                          |
      | HEAD        | git rebase --continue            |
      | main        | git push                         |
      | main        | git checkout -b new_feature main |
      | new_feature | git stash pop                    |
    And I end up on the "new_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILES            |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new_feature | local            | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH      | FILES            | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new_feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git hack --continue `
    Then it runs the Git commands
      | BRANCH      | COMMAND                          |
      | main        | git push                         |
      | main        | git checkout -b new_feature main |
      | new_feature | git stash pop                    |
    And I end up on the "new_feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH      | LOCATION         | MESSAGE                   | FILES            |
      | main        | local and remote | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
      | new_feature | local            | conflicting remote commit | conflicting_file |
      |             |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH      | FILES            | CONTENT          |
      | main        | conflicting_file | resolved content |
      | new_feature | conflicting_file | resolved content |
