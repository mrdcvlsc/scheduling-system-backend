# TODO List

1. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 41

    make a function to compare two `EncodingResource`.

2. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 459

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

3. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 471

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

4. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 699

    analyze if this is really needed?

5. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 769

    [implement below] search for general rooms that are available (consult first)

6. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 771

    [implement below] search available lab room for lecture subjects (consult first)

7. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 929

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

8. File: [./Routes/RoutesV1/GetClassSchedule.go](./Routes/RoutesV1/GetClassSchedule.go), Line: 39

    use department_id for authentication later on.

9. File: [./Routes/RoutesV1/GetClassSchedule.go](./Routes/RoutesV1/GetClassSchedule.go), Line: 185

    use department_id for authentication later on.

10. File: [./Routes/RoutesV1/GenerateSchedule.go](./Routes/RoutesV1/GenerateSchedule.go), Line: 40

    this is just for testing

