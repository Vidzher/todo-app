create unique index if not exists tasks_author_id_title_idx
    on todoapp.tasks(author_user_id, title);