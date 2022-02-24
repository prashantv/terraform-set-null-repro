# Terraform Set Null Repro

This is a repro for what seems to be a terraform bug where an update to
an object in a single-element set somehow leads to an Update call where
the set has multiple elements (even though the old/new state only had a
single element in the set).

The repro provider stores state in a "db" file to keep the repro simple.
The issue was originally found when a backing API rejected updates as the
attributes on the object were invalid (they were empty). The stub provider
simply panics if it receives multiple elements in the job set, since the
schema sets MaxItems: 1.

To reproduce, check out this repo, cd into it, and run the following commands:
```
# ensure the repo is in a clean state
$ git clean -fdx

# install the provider to the local plugins directory
$ make install

# apply the initial state
$ make apply

# MANUAL: update the namespace field to something else in repro/main.tf

# apply will now panic, as the "job" set has 2 elements
$ make apply
[...]
panic: jobSet has more elements than expected: [map[namespace: search_tags:map[]] map[namespace:NS2 search_tags:map[foo:bar]]]
```

This is the plan for the update:
```
  # tftest_service.s will be updated in-place
  ~ resource "tftest_service" "s" {
        id   = "1645686922508288000"
        name = "svc1"

      - job {
          - namespace   = "NS1" -> null
          - search_tags = {
              - "foo" = "bar"
            } -> null
        }
      + job {
          + namespace   = "NS2"
          + search_tags = {
              + "foo" = "bar"
            }
        }
    }

```

We expect an update that only contains the new job, NS2, but we end up with an update
that contains 2 elements, the first is an empty element, and the latter is the correct
element.

Some interesting findings:
* If TypeList is used instead of TypeSet, this issue does not occur.
* If the "job" field does not contain search_tags, then this issue does not occur.

