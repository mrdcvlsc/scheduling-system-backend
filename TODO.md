# TODO List

1. File: [./main.go](./main.go), Line: 52

    implement mongo db persistence implementation

2. File: [./GeneticAlgorithm/Selection.go](./GeneticAlgorithm/Selection.go), Line: 3

    tournament selection

3. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 15

    implement the whole genetic algorithm function

4. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 78

    crossover some individuals and add them to the CURRENT population

5. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 80

    apply random mutations to some of the CURRENT individuals in the population

6. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 82

    apply tournament selection to be saved to the next generation

7. File: [./GeneticAlgorithm/GeneticAlgorithm.go](./GeneticAlgorithm/GeneticAlgorithm.go), Line: 98

    sort individuals base on fitness

8. File: [./GeneticAlgorithm/Crossover.go](./GeneticAlgorithm/Crossover.go), Line: 3

    implement a crossover function

9. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 24

    complete this empty schedule fitness test

10. File: [./GeneticAlgorithm/FitnessFunction_test.go](./GeneticAlgorithm/FitnessFunction_test.go), Line: 103

    complete this generated schedule fitness test

11. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 27

    random subject time slot nudge

12. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 29

    random subject time slot day swap

13. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 64

    random subject clear

14. File: [./GeneticAlgorithm/RandomMutation.go](./GeneticAlgorithm/RandomMutation.go), Line: 65

    then generate/complete encoding

15. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 146

    implement distribution types  front compressed distribution : start  for now the this is a brute force implementation, yet it is still enough considering the current low numbers of courses, instructors, rooms and sections. but I believe this still can be optimize in the future if we wanted to.  or just create another implementation for inidividual generation.

16. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 158

    the map below was added for debugging purposes only, remove when the code becomes stable and optimized.

17. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 455

    [implement below] search for general rooms that are available (consult first)

18. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 457

    [implement below] search available lab room for lecture subjects (consult first)

19. File: [./GeneticAlgorithm/GenerateIndividual.go](./GeneticAlgorithm/GenerateIndividual.go), Line: 588

    when generating solutions while the genetic algorithm is running, we should also generate an index file to be use for querying each sections in the generated university schedules, make the generated schedule and index global for access.

20. File: [./GeneticAlgorithm/FitnessFunction.go](./GeneticAlgorithm/FitnessFunction.go), Line: 92

    implement preference heat map comparison based fitness function

21. File: [./Routes/RoutesV2/ScheduleDelete.go](./Routes/RoutesV2/ScheduleDelete.go), Line: 54

    use department_id for authentication later on.

22. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 34

    use department_id for authentication later on.

23. File: [./Routes/RoutesV1/ScheduleDelete.go](./Routes/RoutesV1/ScheduleDelete.go), Line: 103

    use department_id for authentication later on.

