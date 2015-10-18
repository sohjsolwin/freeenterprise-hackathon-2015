using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace HackSite.Dto
{
    public class UnemploymentData
    {
        public string State { get; set; }
        public string StateAbbr { get; set; }
        public DateTime? Date { get; set; }
        public decimal? Value { get; set; }
    }

    //public class StateData
    //{
    //    public DateTime Date { get; set; }
    //    public double Value { get; set; }
    //}

}