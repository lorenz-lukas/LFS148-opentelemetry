package io.novatec.todobackend;

import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

public interface TodoRepository extends JpaRepository<Todo, Long> {
	List<Todo> findByDoneFalseOrderByCreatedAtAsc();
	Optional<Todo> findFirstByTitleAndDoneFalseOrderByCreatedAtDesc(String title);
}
