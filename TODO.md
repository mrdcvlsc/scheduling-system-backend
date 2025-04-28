# TODO List

1. File: [./GeneticAlgorithm/Selection.go](./GeneticAlgorithm/Selection.go), Line: 3

    tournament selection

2. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 27

    implement the whole genetic algorithm function

3. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 306

    optimization - I think there's no need to set/unset here, we can optimize by direct checking from instructor and room encoding resource availability

4. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 353

    optimization - remove after implementing direct check using instructor and room availability

5. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 364

    optimization - remove after implementing direct check using instructor and room availability

6. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 31

    complete this empty schedule fitness test

7. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 110

    complete this generated schedule fitness test

8. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 362

    debug there is an error here (currently this mutation function is not being used)

9. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 146

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

10. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 158

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

11. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 489

    [implement below] search for general rooms that are available (consult first)

12. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 491

    [implement below] search available lab room for lecture subjects (consult first)

13. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 654

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

14. File: [./GeneticAlgorithm/FitnessFunction.go](./GeneticAlgorithm/FitnessFunction.go), Line: 121

    implement preference heat map comparison based fitness function

15. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

16. File: [./Routes/RoutesV1/SchedulePost.go](./Routes/RoutesV1/SchedulePost.go), Line: 266

    on genetic algorithm error - just use normal schedule generation result

17. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 40

    use department_id for authentication later on.

18. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 115

    use department_id for authentication later on.

