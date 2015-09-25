package services

// Travis

type TRRepository struct {
    Id         int
    Name       string
    Owner_name string
    Url        string
}

type TRPayload struct {
    Author_email        string
    Author_name         string
    Branch              string
    Build_url           string
    Commit              string
    Committed_at        string
    Committer_email     string
    Committer_name      string
    Compare_url         string
    Duration            int
    Finished_at         string
    Id                  int
    Message             string
    Number              string
    Repository          TRRepository
    Started_at          string
    State               string
    Status              int
    Status_message      string
    Type                string
    Pull_request        bool
    Pull_request_number string
    Pull_request_type   string
    Tag                 string
}
