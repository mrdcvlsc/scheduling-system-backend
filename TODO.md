# TODO List

1. File: [./main.go](./main.go), Line: 52

    implement mongo db persistence implementation

2. File: [./GeneticAlgorithm/Selection.go](./GeneticAlgorithm/Selection.go), Line: 3

    tournament selection

3. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 3

    implement the whole genetic algorithm function

4. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 3

    implement a crossover function

5. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 3

    implement random mutation function

6. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 141

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

7. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 153

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

8. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 436

    [implement below] search for general rooms that are available (consult first)

9. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 438

    [implement below] search available lab room for lecture subjects (consult first)

10. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 571

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

11. File: [./GeneticAlgorithm/FitnessFunction.go](./GeneticAlgorithm/FitnessFunction.go), Line: 3

    implement a fitness function

12. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

13. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 34

    use department_id for authentication later on.

14. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 103

    use department_id for authentication later on.

