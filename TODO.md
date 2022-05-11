1. Delete method maybe doesn't do anything at all
2. when creating a single avenger resource, we wind up including *all* resources in our state.
3. #2 happens on *subsequent* apply, not on initial supply which is weird.