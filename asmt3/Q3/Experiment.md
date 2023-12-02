# 3. (10 marks) Discuss whether your fault tolerant version of Ivy still preserves sequential consistency or not.
Sequential consistency is a strong safety property for concurrent systems. Informally, sequential consistency implies that operations appear to take place in some total order, and that that order is consistent with the order of operations on each individual process.

# Experiments 
(Measure the end-to-end time performance with multiple read/write req.)

1. Without any faults, compare the performance of the basic version of Ivy protocol and the new fault tolerant version using requests from at least 10 clients. 
2. Evaluate the new design in the presence of a single fault – one CM fails only once. Specifically, you can simulate two scenarios 
- a) when the primary CM fails at a random time, 
- b) when the primary CM restarts after the failure. 
Compare the performance of these two cases with the equivalent scenarios without any CM faults. 
3. Evaluate the new design in the presence of multiple faults for primary CM – primary CM fails and restarts multiple times. Compare the performance with respect to the equivalent scenarios without any CM faults.
4. Evaluate the new design in the presence of multiple faults for primary CM and backup CM – both primary CM and backup CM fail and restart multiple times. Compare the performance with respect to the equivalent scenarios without any CM faults.