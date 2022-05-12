1. ~Delete method maybe doesn't do anything at all~ Fixed
2. ~when creating a single avenger resource, we wind up including *all* resources in our state.~ Fixed
3. ~#2 happens on *subsequent* apply, not on initial supply which is weird.~
4. Update to change the name definitely doesn't work