# TODO List

1. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 141

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

2. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 153

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

3. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 436

    [implement below] search for general rooms that are available (consult first)

4. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 438

    [implement below] search available lab room for lecture subjects (consult first)

5. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 571

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

6. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

7. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 34

    use department_id for authentication later on.

8. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 103

    use department_id for authentication later on.

