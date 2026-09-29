# Reader saves position after each page turn

**Status:** Proposed

**Date:** 2032-03-18

## Decision

Reader will save the current page on the device after each page turn. Opening the same book will restore the last saved page.

The reading screen will own the current position. The local library will own storage and recovery.

## Context

Reader currently opens every book at the beginning. Readers must find their previous page whenever they reopen a book.

The first version supports one device. Account sync requires a separate design.

## Position storage

The local library stores one position for each book identity. A filename change does not create another position.

After a page turn, the reading screen sends the book identity and page number to the local library. The library replaces the saved position in one transaction.

## Recovery

A failed save does not change the previous position. Reader keeps the current page visible and retries the save.

After an app restart, Reader opens the last committed position. Pages viewed after a failed save may need to be revisited.

## Consequences

Readers can resume without searching for their page. Saving after each page turn adds local writes.

The saved position does not follow the reader to another device.

## Alternatives considered

### Save only when the book closes

This reduces writes but can lose the entire session's progress after a crash.

### Sync through an account

Sync would support several devices but would also require a policy for competing positions. The first version uses local storage.
