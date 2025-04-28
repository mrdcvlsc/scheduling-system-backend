# TODO List

1. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 306

    optimization - I think there's no need to set/unset here, we can optimize by direct checking from instructor and room encoding resource availability

2. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 353

    optimization - remove after implementing direct check using instructor and room availability

3. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 364

    optimization - remove after implementing direct check using instructor and room availability

4. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 362

    debug there is an error here (currently this mutation function is not being used)

5. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 146

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

6. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 158

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

7. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 489

    [implement below] search for general rooms that are available (consult first)

8. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 491

    [implement below] search available lab room for lecture subjects (consult first)

9. File: [./GeneticAlgorithm/FitnessFunction.go](./GeneticAlgorithm/FitnessFunction.go), Line: 124

    implement preference heat map comparison based fitness function

10. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

11. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 40

    use department_id for authentication later on.

12. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 115

    use department_id for authentication later on.

