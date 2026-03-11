package io.novatec.todobackend;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

import io.opentelemetry.instrumentation.annotations.SpanAttribute;
import io.opentelemetry.instrumentation.annotations.WithSpan;

@RestController
@CrossOrigin(origins = "*")
public class TodoController {

	private final Logger logger = LoggerFactory.getLogger(TodoController.class);

	@Value("${HOSTNAME:not_set}")
	String hostname;

	@Autowired
	TodoRepository todoRepository;

	private String getInstanceId() {
		if (!hostname.equals("not_set")) {
			return hostname;
		}
		return "probably localhost";
	}

	@WithSpan
	@GetMapping("/hello")
	String hello() {
		return getInstanceId() + " Hallo, Welt ! ";
	}

	@WithSpan
	@GetMapping("/fail")
	String fail() {
		System.exit(1);
		return "fixed!";
	}

	@WithSpan
	@GetMapping("/todos/")
	List<String> getTodos() {
		List<String> todos = new ArrayList<String>();

		todoRepository.findByDoneFalseOrderByCreatedAtAsc().forEach(todo -> todos.add(todo.getTitle()));
		logger.info("GET /todos/ {}", todos);

		return todos;
	}

	@WithSpan
	@PostMapping("/todos/{todo}")
	String addTodo(@PathVariable String todo) {
		this.someInternalMethod(todo);
		logger.info("POST /todos/ {}", todo);

		return todo;
	}

	@WithSpan
	String someInternalMethod(@SpanAttribute String todo) {
		Todo todoEntity = todoRepository.findFirstByTitleAndDoneFalseOrderByCreatedAtDesc(todo)
			.orElseGet(() -> new Todo(todo, ""));

		todoEntity.setDone(false);
		todoRepository.save(todoEntity);
		if (todo.equals("slow")) {
			try {
				Thread.sleep(1000);
			} catch (InterruptedException e) {
				Thread.currentThread().interrupt();
			}
		}
		if (todo.equals("fail")) {
			System.out.println("Failing ...");
			throw new RuntimeException();
		}
		return todo;
	}

	@WithSpan
	@DeleteMapping("/todos/{todo}")
	String removeTodo(@PathVariable String todo) {
		Optional<Todo> existingTodo = todoRepository.findFirstByTitleAndDoneFalseOrderByCreatedAtDesc(todo);
		if (existingTodo.isEmpty()) {
			logger.info("DELETE /todos/ {} not found", todo);
			return "not found " + todo;
		}

		Todo todoEntity = existingTodo.get();
		todoEntity.setDone(true);
		todoRepository.save(todoEntity);
		logger.info("DELETE /todos/ {}", todo);
		return "done " + todo;
	}
}
