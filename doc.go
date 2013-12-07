/*
Trie is a prefix index package for golang.

Terminology used

Trie - Contains the Root of the Index - which is a Branch

Branch - A Branch might have a LeafValue or not and might have other Branches splitting off. It has a flag `End` that marks the end of a term that has been inserted.

Entry - an entry refers to a _complete_ term that is inserted, removed from, or matched in the index. It requires `End` on the Branch to be set to `true`, which makes it different from a

Prefix - which does not require the Branch to have End set to `true` to match.
*/
package trie
