CREATE TABLE users(
    id integer not null GENERATED ALWAYS AS IDENTITY,
    full_name varchar not null,
	PRIMARY KEY(id)
);
  
   
CREATE TABLE segments(
    id integer not null GENERATED ALWAYS AS IDENTITY,
    title varchar not null unique,
status bool,
	PRIMARY KEY(id)

);

CREATE TABLE actions(
    id integer not null GENERATED ALWAYS AS IDENTITY,
	user_id integer,
	segment_id integer,
	start_date timestamp,
	end_date timestamp,
	PRIMARY KEY(id),
	CONSTRAINT fk_user
   		FOREIGN KEY(user_id) 
      		REFERENCES users(id),
	CONSTRAINT fk_segment
   		FOREIGN KEY(segment_id) 
      		REFERENCES segments(id),
	UNIQUE(user_id, segment_id, start_date)
)
