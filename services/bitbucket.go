package services

// Bitbucket

type BBRepository struct {
    Name       string
    Is_private bool
}

type BBAuthor struct {
    User User
    Raw  string
}

type BBCommits struct {
    Message string
    Hash    string
    Author  BBAuthor
}

type BBChanges struct {
    Commits []BBCommits
}

type BBPush struct {
    Changes []BBChanges
}

type BBPayload struct {
    Repository BBRepository
    Push       BBPush
}
